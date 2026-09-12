import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import process from "node:process";
import { z } from "zod";

const server = new McpServer({
    name: "ts-greeter",
    version: "1.0.0",
});

server.tool(
    "greet",
    "say hi",
    { name: z.string().describe("the name of the person to greet") },
    async ({ name }) => {
        return {
            content: [{ type: "text", text: `Hi my friend ${name}` }],
        };
    }
);

async function main() {
    const transport = new StdioServerTransport();
    await server.connect(transport);
}

main().catch((error) => {
    console.error("Server error:", error);
    process.exit(1);
});
