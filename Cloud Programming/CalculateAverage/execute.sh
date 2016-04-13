# Do not uncomment these lines to directly execute the script
# Modify the path to fit your need before using this script
#hdfs dfs -rm -r /user/TA/CalculateAverage/Output/
#hadoop jar CalculateAverage.jar calculateAverage.CalculateAverage /user/shared/CalculateAverage/Input /user/TA/CalculateAverage/Output
#hdfs dfs -cat /user/TA/CalculateAverage/Output/part-*

hdfs dfs -rm -r /user/s101062115/CalculateAverage/output/
hadoop jar CalculateAverage.jar calculateAverage.CalculateAverage /user/s101062115/input2 /user/s101062115/CalculateAverage/output
hdfs dfs -cat /user/s101062115/CalculateAverage/output/part-*
