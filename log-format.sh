#!/bin/bash

bash -c "GREEN='\033[0;32m'; ORANGE='\033[0;33m'; RED='\033[0;31m'; NC='\033[0m';
while IFS= read -r line; do
        if [[ \$line == *'level=INFO'* ]]; then
                msg=\$(echo \"\$line\" | grep -oP 'msg=\"\\K[^\"]+');
                echo -e \"\${GREEN}\${msg}\${NC}\";
        elif [[ \$line == *'level=WARN'* ]]; then
                msg=\$(echo \"\$line\" | grep -oP 'msg=\"\\K[^\"]+');
                echo -e \"\${ORANGE}\${msg}\${NC}\";
        elif [[ \$line == *'level=ERROR'* ]]; then
                msg=\$(echo \"\$line\" | grep -oP 'msg=\"\\K[^\"]+');
                echo -e \"\${RED}\${msg}\${NC}\";
        fi;
done"
