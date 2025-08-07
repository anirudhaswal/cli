import {
  quicktype,
  JSONSchemaInput,
  FetchingJSONSchemaStore,
  InputData,
  JSONInput,
} from "npm:quicktype-core@23.2.6";
import $RefParser from "npm:@apidevtools/json-schema-ref-parser";


async function quicktypeJSONSchema(
  targetLanguage: any,
  typeName: string,
  jsonSchemaString: string,
  rendererOptions: Record<string, string> = {},
) {
  const schemaInput = new JSONSchemaInput(new FetchingJSONSchemaStore());
   
  let flattenedSchema;
  try {
    const dereferencedSchema = await $RefParser.dereference(JSON.parse(jsonSchemaString), {
      dereference: { circular: "ignore" },
    });
    flattenedSchema = dereferencedSchema
  } catch (error) {
    console.error("Error parsing or dereferencing schema:", error);
    throw new Error("Invalid JSON schema or failed to resolve references");
  }

  await schemaInput.addSource({
    name: typeName,
    schema: JSON.stringify(flattenedSchema),
  });

  const inputData = new InputData();
  inputData.addInput(schemaInput);

  const renderOptions = {
    "just-types": "true",
    "no-extra-properties": "true", 
    "no-optional-null": "true",
    ...rendererOptions
  };


  return await quicktype({
    inputData,
    lang: targetLanguage,
    rendererOptions: renderOptions
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
  let existingContent = "";
  try {
    existingContent = await Deno.readTextFile(outputPath);
  } catch(_) {}

  const finalOutput = existingContent 
    ? existingContent + "\n\n" + output 
    : output;

  await Deno.writeTextFile(outputPath, finalOutput);

}

main();