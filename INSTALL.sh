[ "$UID" -eq 0 ] || exec sudo bash "$0" "$SHELL" "$HOME"
set -e
INSTALL_ROOT="/usr/local/bin"
PARENT_SHELL=$1
PARENT_HOME=$2


if [[ "$PARENT_SHELL" == *"zsh"* ]]; then
  RC_FILE="$PARENT_HOME/.zshrc"
  SCRIPT_EXTENSION="zsh"
else
  RC_FILE="$PARENT_HOME/.bashrc"
  SCRIPT_EXTENSION="sh"
fi

echo "Installing core..."
cd gimme-core && go install gimme-core.go && cd ..
echo "Done."

echo "Installing scripts..."
cp ./scripts/gimme.$SCRIPT_EXTENSION "$INSTALL_ROOT"
echo "Done."

echo "Installing alias..."

grep -s "alias gimme" "$RC_FILE"
if [[ $? == 1 ]]; then
  echo "alias gimme='. gimme.$SCRIPT_EXTENSION'" >> "$RC_FILE"
fi

echo -e "Done. \nRun the following to begin using gimme:\nsource $RC_FILE"