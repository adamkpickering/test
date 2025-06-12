package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types/ref"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	host := config.Host{
		Name: "registry.suse.com/rancher",
		User: os.Getenv("PRIME_USERNAME"),
		Pass: os.Getenv("PRIME_PASSWORD"),
	}
	rc := regclient.New(regclient.WithConfigHost(host))
	ref, err := ref.New("registry.suse.com/rancher/mirrored-kubernetes-external-dns:v0.7.3")
	if err != nil {
		return fmt.Errorf("failed to create ref: %w", err)
	}
	defer rc.Close(context.Background(), ref)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	man, err := rc.ManifestHead(ctx, ref)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to manifest head: %w", err)
	}
	defer cancel()

	fmt.Printf("got manifest: %v\n", man)
	return nil
}

func (opts *rootOpts) processRef(ctx context.Context, s ConfigSync, src, tgt ref.Ref, action actionType) error {
	mSrc, err := opts.rc.ManifestHead(ctx, src, regclient.WithManifestRequireDigest())
	if err != nil && errors.Is(err, errs.ErrUnsupportedAPI) {
		mSrc, err = opts.rc.ManifestGet(ctx, src)
	}
	if err != nil {
		slog.Error("Failed to lookup source manifest",
			slog.String("source", src.CommonName()),
			slog.String("error", err.Error()))
		return err
	}
	fastCheck := (s.FastCheck != nil && *s.FastCheck)
	forceRecursive := (s.ForceRecursive != nil && *s.ForceRecursive)
	referrers := (s.Referrers != nil && *s.Referrers)
	digestTags := (s.DigestTags != nil && *s.DigestTags)
	mTgt, err := opts.rc.ManifestHead(ctx, tgt, regclient.WithManifestRequireDigest())
	tgtExists := (err == nil)
	tgtMatches := false
	if err == nil && manifest.GetDigest(mSrc).String() == manifest.GetDigest(mTgt).String() {
		tgtMatches = true
	}
	if tgtMatches && (fastCheck || (!forceRecursive && !referrers && !digestTags)) {
		slog.Debug("Image matches",
			slog.String("source", src.CommonName()),
			slog.String("target", tgt.CommonName()))
		return nil
	}
	if tgtExists && action == actionMissing {
		slog.Debug("target exists",
			slog.String("source", src.CommonName()),
			slog.String("target", tgt.CommonName()))
		return nil
	}

	// skip when source manifest is an unsupported type
	smt := manifest.GetMediaType(mSrc)
	if !slices.Contains(s.MediaTypes, smt) {
		slog.Info("Skipping unsupported media type",
			slog.String("ref", src.CommonName()),
			slog.String("mediaType", manifest.GetMediaType(mSrc)),
			slog.Any("allowed", s.MediaTypes))
		return nil
	}

	// if platform is defined and source is a list, resolve the source platform
	if mSrc.IsList() && s.Platform != "" {
		platDigest, err := opts.getPlatformDigest(ctx, src, s.Platform, mSrc)
		if err != nil {
			return err
		}
		src.Digest = platDigest.String()
		if tgtExists && platDigest.String() == manifest.GetDigest(mTgt).String() {
			tgtMatches = true
		}
		if tgtMatches && (s.ForceRecursive == nil || !*s.ForceRecursive) {
			slog.Debug("Image matches for platform",
				slog.String("source", src.CommonName()),
				slog.String("platform", s.Platform),
				slog.String("target", tgt.CommonName()))
			return nil
		}
	}
	if tgtMatches {
		slog.Info("Image refreshing",
			slog.String("source", src.CommonName()),
			slog.String("target", tgt.CommonName()),
			slog.Bool("forced", forceRecursive),
			slog.Bool("digestTags", digestTags),
			slog.Bool("referrers", referrers))
	} else {
		slog.Info("Image sync needed",
			slog.String("source", src.CommonName()),
			slog.String("target", tgt.CommonName()))
	}
	if action == actionCheck {
		return nil
	}

	// wait for parallel tasks
	throttleDone, err := opts.throttle.Acquire(ctx, throttle{})
	if err != nil {
		return fmt.Errorf("failed to acquire throttle: %w", err)
	}
	// delay for rate limit on source
	if s.RateLimit.Min > 0 && manifest.GetRateLimit(mSrc).Set {
		// refresh current rate limit after acquiring throttle
		mSrc, err = opts.rc.ManifestHead(ctx, src)
		if err != nil {
			slog.Error("rate limit check failed",
				slog.String("source", src.CommonName()),
				slog.String("error", err.Error()))
			throttleDone()
			return err
		}
		// delay if rate limit exceeded
		rlSrc := manifest.GetRateLimit(mSrc)
		for rlSrc.Remain < s.RateLimit.Min {
			throttleDone()
			slog.Info("Delaying for rate limit",
				slog.String("source", src.CommonName()),
				slog.Int("source-remain", rlSrc.Remain),
				slog.Int("source-limit", rlSrc.Limit),
				slog.Int("step-min", s.RateLimit.Min),
				slog.Duration("sleep", s.RateLimit.Retry))
			select {
			case <-ctx.Done():
				return ErrCanceled
			case <-time.After(s.RateLimit.Retry):
			}
			throttleDone, err = opts.throttle.Acquire(ctx, throttle{})
			if err != nil {
				return fmt.Errorf("failed to reacquire throttle: %w", err)
			}
			mSrc, err = opts.rc.ManifestHead(ctx, src)
			if err != nil {
				slog.Error("rate limit check failed",
					slog.String("source", src.CommonName()),
					slog.String("error", err.Error()))
				throttleDone()
				return err
			}
			rlSrc = manifest.GetRateLimit(mSrc)
		}
		slog.Debug("Rate limit passed",
			slog.String("source", src.CommonName()),
			slog.Int("source-remain", rlSrc.Remain),
			slog.Int("step-min", s.RateLimit.Min))
	}
	defer throttleDone()

	// verify context has not been canceled while waiting for throttle
	select {
	case <-ctx.Done():
		return ErrCanceled
	default:
	}

	// run backup
	if tgtExists && !tgtMatches && s.Backup != "" {
		// expand template
		data := struct {
			Ref  ref.Ref
			Step ConfigSync
			Sync ConfigSync
		}{Ref: tgt, Step: s, Sync: s}
		backupStr, err := template.String(s.Backup, data)
		if err != nil {
			slog.Error("Failed to expand backup template",
				slog.String("original", tgt.CommonName()),
				slog.String("backup-template", s.Backup),
				slog.String("error", err.Error()))
			return err
		}
		backupStr = strings.TrimSpace(backupStr)
		backupRef := tgt
		if strings.ContainsAny(backupStr, ":/") {
			// if the : or / are in the string, parse it as a full reference
			backupRef, err = ref.New(backupStr)
			if err != nil {
				slog.Error("Failed to parse backup reference",
					slog.String("original", tgt.CommonName()),
					slog.String("template", s.Backup),
					slog.String("backup", backupStr),
					slog.String("error", err.Error()))
				return err
			}
		} else {
			// else parse backup string as just a tag
			backupRef = backupRef.SetTag(backupStr)
		}
		defer opts.rc.Close(ctx, backupRef)
		// run copy from tgt ref to backup ref
		slog.Info("Saving backup",
			slog.String("original", tgt.CommonName()),
			slog.String("backup", backupRef.CommonName()))
		err = opts.rc.ImageCopy(ctx, tgt, backupRef)
		if err != nil {
			// Possible registry corruption with existing image, only warn and continue/overwrite
			slog.Warn("Failed to backup existing image",
				slog.String("original", tgt.CommonName()),
				slog.String("template", s.Backup),
				slog.String("backup", backupRef.CommonName()),
				slog.String("error", err.Error()))
		}
	}

	rcOpts := []regclient.ImageOpts{}
	if s.DigestTags != nil && *s.DigestTags {
		rcOpts = append(rcOpts, regclient.ImageWithDigestTags())
	}
	if s.Referrers != nil && *s.Referrers {
		if len(s.ReferrerFilters) == 0 {
			rcOpts = append(rcOpts, regclient.ImageWithReferrers())
		} else {
			for _, filter := range s.ReferrerFilters {
				rOpts := []scheme.ReferrerOpts{}
				if filter.ArtifactType != "" {
					rOpts = append(rOpts, scheme.WithReferrerMatchOpt(descriptor.MatchOpt{ArtifactType: filter.ArtifactType}))
				}
				if filter.Annotations != nil {
					rOpts = append(rOpts, scheme.WithReferrerMatchOpt(descriptor.MatchOpt{Annotations: filter.Annotations}))
				}
				rcOpts = append(rcOpts, regclient.ImageWithReferrers(rOpts...))
			}
		}
		if s.ReferrerSrc != "" {
			referrerSrc, err := ref.New(s.ReferrerSrc)
			if err != nil {
				slog.Error("failed to parse referrer source reference",
					slog.String("referrerSource", s.ReferrerSrc),
					slog.String("error", err.Error()))
			}
			rcOpts = append(rcOpts, regclient.ImageWithReferrerSrc(referrerSrc))
		}
		if s.ReferrerTgt != "" {
			referrerTgt, err := ref.New(s.ReferrerTgt)
			if err != nil {
				slog.Error("failed to parse referrer target reference",
					slog.String("referrerTarget", s.ReferrerTgt),
					slog.String("error", err.Error()))
			}
			rcOpts = append(rcOpts, regclient.ImageWithReferrerTgt(referrerTgt))
		}
	}
	if s.FastCheck != nil && *s.FastCheck {
		rcOpts = append(rcOpts, regclient.ImageWithFastCheck())
	}
	if s.ForceRecursive != nil && *s.ForceRecursive {
		rcOpts = append(rcOpts, regclient.ImageWithForceRecursive())
	}
	if s.IncludeExternal != nil && *s.IncludeExternal {
		rcOpts = append(rcOpts, regclient.ImageWithIncludeExternal())
	}
	if len(s.Platforms) > 0 {
		rcOpts = append(rcOpts, regclient.ImageWithPlatforms(s.Platforms))
	}

	// Copy the image
	slog.Debug("Image sync running",
		slog.String("source", src.CommonName()),
		slog.String("target", tgt.CommonName()))
	err = opts.rc.ImageCopy(ctx, src, tgt, rcOpts...)
	if err != nil {
		slog.Error("Failed to copy image",
			slog.String("source", src.CommonName()),
			slog.String("target", tgt.CommonName()),
			slog.String("error", err.Error()))
		return err
	}
	return nil
}
