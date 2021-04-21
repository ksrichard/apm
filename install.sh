if ! command -v curl &> /dev/null
then
    echo "'curl' could not be found, please install it!"
    exit
fi

if ! command -v jq &> /dev/null
then
    echo "'jq' could not be found, please install it!"
    exit
fi

os=""
target=""
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        os="linux"
        target="/usr/local/bin/apm"
elif [[ "$OSTYPE" == "darwin"* ]]; then
        os="darwin"
        target="/usr/local/bin/apm"
else
        echo "Not supported OS!"
        exit 1
fi
echo "Installing/Updating apm to ${target}..."
rm -rf ${target}
curl -L $(curl -s https://api.github.com/repos/ksrichard/apm/releases/latest | jq -r ".assets[] | select(.name | test(\"${os}_amd64\")) | .browser_download_url") --output ${target}
chmod +x ${target}
echo "apm is ready to use!"