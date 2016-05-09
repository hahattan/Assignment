package calculateAverage;

import java.io.IOException;
import java.util.StringTokenizer;

import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

public class CalculateAverageMapper extends Mapper<LongWritable, Text, Text, SumCountPair> {

  private SumCountPair one = new SumCountPair();
  private Text first = new Text();

  public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

    // we simply use StringTokenizer to split the words for us.
    StringTokenizer itr = new StringTokenizer(value.toString());
    while (itr.hasMoreTokens()) {

      String word = itr.nextToken();
      String num = itr.nextToken();
      first.set(word);
      one = new SumCountPair(Integer.parseInt(num), 1);

      // create <K, V> pair
      context.write(first, one);

    }

  }

}
