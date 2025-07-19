#!/bin/bash
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 old_name new_name"
    exit 1
fi

old_name=$1
new_name=$2

# Update go.mod module name
sed -i "s|module $old_name|module $new_name|" go.mod

# Update imports in all .go and .templ files
find . -type f \( -name "*.go" -o -name "*.templ" \) -exec sed -i "s|\"$old_name/|\"$new_name/|g" {} +
