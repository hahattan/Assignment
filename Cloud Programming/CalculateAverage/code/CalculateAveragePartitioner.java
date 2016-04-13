package calculateAverage;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Partitioner;

public class CalculateAveragePartitioner extends Partitioner<Text, SumCountPair> {
	@Override
	public int getPartition(Text key, SumCountPair values, int numReduceTasks) {
		// customize which <K ,V> will go to which reducer
		if((key.charAt(0) >= 65 && key.charAt(0) <= 71)) {
			return 0;
		}else {
			return 1;
		}
	}
}
