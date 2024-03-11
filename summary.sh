#!/bin/bash

OUTPUT="README.md"
OS=(windows darwin linux ios android)

truncate -s 0 $OUTPUT
for i in "${OS[@]}"; do
  echo "## $i" >> $OUTPUT
  echo "\`\`\`" >> $OUTPUT
  echo $(cat $i) >> $OUTPUT
  echo "\`\`\`" >> $OUTPUT
done
