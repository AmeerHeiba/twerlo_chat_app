#!/bin/bash

OUTPUT="merged_go_files.txt"

# Clear existing output file
> "$OUTPUT"

# Find all .go files recursively and process them
find . -type f -name "*.go" | while read -r file; do
    echo "// FILE: $file" >> "$OUTPUT"
    echo "" >> "$OUTPUT"
    cat "$file" >> "$OUTPUT"
    echo -e "\n\n// ===== END OF $file =====\n\n" >> "$OUTPUT"
done

echo "All .go files merged into $OUTPUT"