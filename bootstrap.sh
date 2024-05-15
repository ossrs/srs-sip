add_gopath_to_path() {
    GOPATH=$(go env GOPATH)
    # Check if GOPATH is set
    if [ -z "$GOPATH" ]; then
        echo "GOPATH is not set."
        return 1
    fi

    # Check if $GOPATH/bin is already in ~/.bashrc
    if grep -q "$GOPATH/bin" ~/.bashrc; then
        echo "$GOPATH/bin is already in PATH."
        return 0
    fi

    # Add $GOPATH/bin to PATH
    echo "export PATH=\$PATH:$GOPATH/bin" >> ~/.bashrc
    source ~/.bashrc

    echo "$GOPATH/bin has been added to PATH."
}


if ! command -v mage &> /dev/null
then
  pushd /tmp

  OS_IS_LINUX=$(uname -s |grep -q Linux && echo YES)
  if [ "$OS_IS_LINUX" == "YES" ]; then
    add_gopath_to_path
    if [ $? -eq 1 ]; then
        echo "error: Failed to add $GOPATH/bin to PATH."
        exit 1
    fi
  fi

  git clone https://github.com/magefile/mage
  cd mage
  go run bootstrap.go
  rm -rf /tmp/mage
  popd
fi

if ! command -v mage &> /dev/null
then
  echo "error: Ensure `go env GOPATH`/bin is in your \$PATH"
  exit 1
fi

go mod download