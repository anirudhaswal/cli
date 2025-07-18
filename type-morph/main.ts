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
    schemaInput = "./schema.json",
    schemaName = "SchemaType",
    outputPath = "./output.txt",
  ] = Deno.args;

  let text: string;

  try {
    text = await Deno.readTextFile(schemaInput);
  } catch (_) {
    try {
      JSON.parse(schemaInput); // validate it's valid JSON
      text = schemaInput;
    } catch (err) {
      console.error(err);
      Deno.exit(1);
    }
  }

  const { lines } = await quicktypeJSONSchema(language, schemaName, text);

  const output = lines.join("\n");
  await Deno.writeTextFile(outputPath, output);

  console.log(`Generated ${language} types written to ${outputPath}`);
}

main();