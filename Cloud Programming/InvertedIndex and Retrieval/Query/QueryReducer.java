package query;

import java.io.*;
import java.lang.*;
import java.util.HashMap;

import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapred.Counters;
import org.apache.hadoop.conf.Configuration;

public class QueryReducer extends Reducer<QueryKeyPair,QueryValuePair,Text,Text> {

    private String last = "";
    private Text rank = new Text();
    private Text result = new Text();
    private HashMap<Integer, String> hash = new HashMap<Integer, String>();

    public void setup(Context context) throws IOException, InterruptedException {

        Configuration conf = context.getConfiguration();
        Path file = new Path(conf.get("path"));

        try {
    	    FileSystem fs = FileSystem.get(conf);
    	    Path path = new Path("InvertedIndex/id.log");
    	    BufferedReader br = new BufferedReader(new InputStreamReader(fs.open(path)));
    	    String line;
    	    line = br.readLine();
            line = br.readLine();
            int i = 1;
    	    while(line != null) {
                hash.put(i++, line);
                line = br.readLine();
            }
    	}catch(Exception e) {

    	}
    }

    public void reduce(QueryKeyPair key, Iterable<QueryValuePair> values, Context context)
    throws IOException, InterruptedException {

        long rk = context.getCounter(TOP10_COUNTER.RANK).getValue();
	    for (QueryValuePair val: values) {
		long limit = context.getCounter(TOP10_COUNTER.NUMBER).getValue();
		if(limit >= 10) break;

	        StringBuilder path = new StringBuilder(context.getConfiguration().get("path"));
	        path.append("/");
	        path.append(hash.get(val.getDocID()));

	        StringBuilder temp = new StringBuilder();
	        temp.append("Rank " + String.valueOf(rk+1) + " : " + hash.get(val.getDocID()));
	        temp.append(":\tscore = " + key.getWeight());
            rank.set(temp.toString());

            temp = new StringBuilder();
            temp.append("\n************************\n");
            String offset[] = val.getOffset().split("\\D");
            try {
        	    FileSystem fs = FileSystem.get(context.getConfiguration());
        	    Path file = new Path(path.toString());
                for(int i = 0; i < offset.length; i++) {
			    if(offset[i].length() == 0) continue;
        	        BufferedReader br = new BufferedReader(new InputStreamReader(fs.open(file)));
                    br.skip(Long.parseLong(offset[i]));
                    String line = br.readLine();
                    temp.append(offset[i] + "\t=>    " + line + "\n");
                }
        	}catch(Exception e) {

        	}
            temp.append("************************");
	        result.set(temp.toString());

            context.getCounter(TOP10_COUNTER.RANK).increment(1);
	        context.getCounter(TOP10_COUNTER.NUMBER).increment(1);
	        context.write(rank, result);
	    }
    }
}
