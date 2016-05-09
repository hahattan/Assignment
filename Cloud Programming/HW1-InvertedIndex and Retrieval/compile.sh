rm -r class/*
javac -classpath  $HADOOP_HOME/share/hadoop/mapreduce/hadoop-mapreduce-client-core-2.7.2.jar:$HADOOP_HOME/share/hadoop/common/hadoop-common-2.7.2.jar -d class InvertedIndex/*
jar -cvf InvertedIndex.jar -C class .

javac -classpath  $HADOOP_HOME/share/hadoop/mapreduce/hadoop-mapreduce-client-core-2.7.2.jar:$HADOOP_HOME/share/hadoop/common/hadoop-common-2.7.2.jar -d class Query/*
jar -cvf Query.jar -C class .
