#!/bin/bash

# Determine the directory of the script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Define source and target directories relative to the script's location
SOURCE_DIR="$(realpath "$SCRIPT_DIR/..")"
TARGET_DIR="$(realpath "$SCRIPT_DIR/../../openk_flat")"

# Step 1: Check if the target directory exists, and remove it if it does
if [ -d "$TARGET_DIR" ]; then
  echo "Target directory exists. Removing it: $TARGET_DIR"
  rm -rf "$TARGET_DIR"
fi

# Step 2: Create the target directory
echo "Creating target directory: $TARGET_DIR"
mkdir -p "$TARGET_DIR"

# Step 3: Generate repo_dir.txt
REPO_DIR_FILE="$SOURCE_DIR/repo_dir.txt"
echo "Generating repo_dir.txt at $REPO_DIR_FILE"
tree --dirsfirst "$SOURCE_DIR" > "$REPO_DIR_FILE"
echo "repo_dir.txt generated."

# Step 4: Function to link files
link_file() {
  local src_file="$1"
  local base_name
  base_name=$(basename "$src_file")
  ln -s "$src_file" "$TARGET_DIR/$base_name"
  echo "Linked: $src_file -> $TARGET_DIR/$base_name"
}

# Step 5: Link repo_dir.txt as the first file
link_file "$REPO_DIR_FILE"

# Step 6: Start flatification
echo "Starting flatification..."

# Link documentation files
for file in "$SOURCE_DIR/docs/project_description.md" \
            "$SOURCE_DIR/docs/shared-vision.md" \
            "$SOURCE_DIR/docs/DEVELOPMENT.md" \
            "$SOURCE_DIR/docs/adr/"* \
            "$SOURCE_DIR/docs/specs/"* \
            "$SOURCE_DIR/docs/models/"*; do
  link_file "$file"
done

# Link core framework files
for file in "$SOURCE_DIR/internal/crypto/"* \
            "$SOURCE_DIR/internal/opene/"* \
            "$SOURCE_DIR/internal/ctx/"* \
            "$SOURCE_DIR/internal/buildinfo/"*; do
  [[ $file == *_test.go ]] && continue
  link_file "$file"
done

# Link logging files
link_file "$SOURCE_DIR/internal/logging/log_error.go"
link_file "$SOURCE_DIR/internal/logging/logger.go"

# Link storage files
for file in "$SOURCE_DIR/internal/storage/store.go"* \
            "$SOURCE_DIR/internal/storage/models/"* \
            "$SOURCE_DIR/internal/storage/memory/"* \
            "$SOURCE_DIR/internal/storage/testsuite/userstore/"*; do
  [[ $file == *_test.go ]] && continue
  link_file "$file"
done

# Link server files
for file in "$SOURCE_DIR/internal/server/interceptors/"* \
            "$SOURCE_DIR/internal/server/config.go" \
            "$SOURCE_DIR/internal/server/grpc_options.go" \
            "$SOURCE_DIR/internal/server/grpc_server.go"; do
  [[ $file == *_test.go ]] && continue
  link_file "$file"
done

# Link all proto files under openk/proto/openk
find "$SOURCE_DIR/proto/openk" -type f -name '*.proto' | while read -r file; do
  link_file "$file"
done

# Link CLI and client implementations
for file in "$SOURCE_DIR/cmd/openk/openk.go" \
            "$SOURCE_DIR/internal/app/client/"* \
            "$SOURCE_DIR/internal/app/app_context.go" \
            "$SOURCE_DIR/internal/app/server.go" \
            "$SOURCE_DIR/internal/cli/cli.go" \
            "$SOURCE_DIR/internal/cli/auth/"* \
            "$SOURCE_DIR/internal/cli/server/"*; do
  [[ $file == *_test.go ]] && continue
  link_file "$file"
done

# Link build and project files
link_file "$SOURCE_DIR/makefile"

# Link all openk workfiles and scripts
for file in "$SOURCE_DIR/scripts/"* \
            "$SOURCE_DIR/openk-"*; do
  link_file "$file"
done

echo "Flatification complete. All files linked to: $TARGET_DIR"
