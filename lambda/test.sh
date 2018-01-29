#!/bin/bash
FUNCTION_NAME="mustafa"
mkdir -p test_results

for i in $(seq 128 64 3008); do
    echo "Setting memory to $i MB"
    aws lambda update-function-configuration --function-name $FUNCTION_NAME --memory-size $i > /dev/null
    for t in {1..4}; do
        for j in {1..10}; do
            filename="test_results2/test_${i}_${t}_${j}"
            echo $filename

            aws lambda invoke \
                --invocation-type RequestResponse \
                --function-name $FUNCTION_NAME \
                --region us-west-2 \
                --cli-read-timeout 300 \
                --payload '{"Iterations":25, "Batch": 100, "Threads": '$t'}' \
                $filename > /dev/null && cat $filename | jq .Mean

            cat $filename | jq ".Durations[]" >> test_results2/durations_${i}_${t}
        done
    done
done
