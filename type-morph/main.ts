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
  });
}

async function main() {
    const [
    language = "typescript",
    schemaPath = "./schema.json",
    outputPath = "./output.txt",
  ] = Deno.args;

  const text = await Deno.readTextFile(schemaPath);
  const { lines } = await quicktypeJSONSchema(language, "SchemaType", text);

  const output = lines.join("\n");
  await Deno.writeTextFile(outputPath, output);

  console.log(`Generated ${language} types written to ${outputPath}`);
}

main();