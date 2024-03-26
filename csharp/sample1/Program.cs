// See https://aka.ms/new-console-template for more information

using HelloWorld;
using System.IO;

var asdf = new MyClass();
asdf.MyNewFunc();
asdf.ReadFile();

namespace HelloWorld
{
	class MyClass
	{

		public void ReadFile()
		{
			string contents = File.ReadAllText("test.txt");
			Console.WriteLine(contents);
		}
		public void MyNewFunc()
		{
			string[] messages = {
				"hello world",
				"volvo",
				"apple",
				"pear",
			};
			foreach (string message in messages)
			{
				if (message == "apple")
				{
					break;
				}
				Console.WriteLine(message);
			}
		}
	}
}
