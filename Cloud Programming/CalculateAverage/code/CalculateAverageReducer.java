package calculateAverage;

import java.io.IOException;

import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

public class CalculateAverageReducer extends Reducer<Text,SumCountPair,Text,DoubleWritable> {

  private DoubleWritable result = new DoubleWritable();

  public void reduce(Text key, Iterable<SumCountPair> values, Context context) throws IOException, InterruptedException {
    int sum = 0;
    int count = 0;
    double avg = 0;
    // agrregate the amount of same starting character
    for (SumCountPair val: values) {
      // count += ...
      sum += val.getSum();
      count += val.getCount();
    }
    avg = (double)sum/(double)count;
    result.set(avg);
    // write the result
    context.write(key, result);

  }
}
