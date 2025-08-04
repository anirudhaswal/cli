import {
  quicktype,
  JSONSchemaInput,
  FetchingJSONSchemaStore,
  InputData,
} from "npm:quicktype-core@23.2.6";

async function quicktypeJSONSchema(
  targetLanguage: any,
  typeName: string,
  jsonSchemaString: string,
  rendererOptions: Record<string, string> = {},
) {
  const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore());

  await schemaInput.addSource({
    name: typeName,
    schema: jsonSchemaString,
  });

  const inputData = new InputData();
  inputData.addInput(schemaInput);

  return await quicktype({
    inputData,
    lang: targetLanguage,
    rendererOptions: rendererOptions
  });
}

function extractRendererOptions(args: string[]): Record<string, string> {
  const options: Record<string, string> = {};
  for (const arg of args) {
    if (arg.startsWith("--build-flags=")) {
      const flags = arg.split("=")[1];
      for (const flag of flags.split(",")) {
        const trimmed = flag.trim();
        if (trimmed) {
          options[trimmed] = "true";
        }
      }
    }
  }
  return options;
}

async function main() {
    const [
    language = "typescript",
    schemaInput = "./schema.json",
    schemaName = "SchemaType",
    outputPath = "./output.txt",
    ...extraArgs
  ] = Deno.args;

  let text: string;

  const rendererOptions = extractRendererOptions(extraArgs);

  try {
    text = await Deno.readTextFile(schemaInput);
  } catch (_) {
    try {
      JSON.parse(schemaInput);
      text = schemaInput;
    } catch (err) {
      console.error(err);
      Deno.exit(1);
    }
  }

  const { lines } = await quicktypeJSONSchema(language, schemaName, text, rendererOptions);

  const output = lines.join("\n");
  await Deno.writeTextFile(outputPath, output);

  console.log(`Generated ${language} types written to ${outputPath}`);
}

main();