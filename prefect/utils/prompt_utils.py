def build_prompt_no_schema(text: str, target: str, known_names: list[str]):
    known_names_formatted = ", ".join(f'"{name}"' for name in known_names)

    return f"""
        You are an expert at structured data extraction. You will be given unstructured text from the provided unstructured text and should convert it into the given structure.
        
        ## Target Person
        The individual whose profile you must extract is: **{target}**

        ## Known Aliases
        The following names MUST be included in the `"names"` array in the final JSON output if they refer to the target person:
        [{known_names_formatted}]
        
        You may also extract additional names or aliases used to refer to this person in the text if they are clearly linked to the target.

        ## Text to Analyze:
        {text}

        ## Task Instructions:
        1. **Extract structured data** based on the schema provided.
        2. **Use the descriptions for reference only.**
        3. **DO NOT return descriptions** in the final output.
        4. **Preserve field structure** exactly as specified.
        5. **Return only extracted values in JSON format.**
        6. **If a field is not found, leave it empty (do not generate missing values or invent information).**
        7. **For lists, return an empty array `[]` if no data is found.**
        8. **Ensure proper data formatting**, e.g.:
        - Dates in ISO 8601 format (`YYYY-MM-DD`).
        - Numeric values as numbers, not strings.
        - Use correct currency codes (`USD`, `EUR`, etc.).
        - Standardized industry names (`Technology`, `Finance`, etc.).
        9. **For nested fields, preserve hierarchy and structure.**
        10. **Exclude any additional commentary or explanations.**

    """
