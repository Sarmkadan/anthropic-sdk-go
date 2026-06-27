using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Anthropic.Sdk;

namespace BasicUsage
{
    class Program
    {
        static void Main(string[] args)
        {
            // Create a new client with your API key
            var client = new AnthropicClient("YOUR_API_KEY");

            // Create a new message
            var message = new Message
            {
                Text = "Hello, world!"
            };

            // Send the message
            var response = client.Messages.New(message);

            // Print the response
            Console.WriteLine(response.Content);
        }
    }
}
