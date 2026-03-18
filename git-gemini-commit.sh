#!/bin/bash

# Get the staged git diff.
GIT_DIFF=$(git diff --staged)

if [ -z "$GIT_DIFF" ]; then
  echo "No staged changes to commit."
  exit 1
fi

# Construct the prompt for Gemini.
# The `gemini` command here is a placeholder for however you would
# run a prompt through the Gemini model in your shell.
PROMPT="Based on the following git diff, please generate a commit message following the Conventional Commits specification. The message should have a type, a scope (optional), and a subject, followed by an optional body.

---
$GIT_DIFF
---
"

# Call Gemini to generate the commit message.
# You will need to have the gemini CLI installed and configured.
GENERATED_MESSAGE=$(gemini "$PROMPT")

# Check if Gemini returned a message.
if [ -z "$GENERATED_MESSAGE" ]; then
  echo "Gemini did not return a commit message."
  exit 1
fi

# Allow the user to edit the generated message using gum.
EDITED_MESSAGE=$(echo "$GENERATED_MESSAGE" | gum write --width 80 --placeholder "Review and edit the commit message")

# If the user cancels the edit, gum might return an empty string.
if [ -z "$EDITED_MESSAGE" ]; then
    echo "Commit cancelled."
    exit 1
fi

# Use the edited message to commit.
echo "$EDITED_MESSAGE" | git commit -F -

echo "Commit created successfully."
