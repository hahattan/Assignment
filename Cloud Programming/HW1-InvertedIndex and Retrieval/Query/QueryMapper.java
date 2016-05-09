package query;

import java.io.*;
import java.util.*;
import java.lang.*;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapred.Counters;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.filecache.DistributedCache;


public class QueryMapper extends Mapper<LongWritable, Text, QueryKeyPair, QueryValuePair> {

    private Double N;
    private Double df;
    private List<String> keyword = new ArrayList<String>();

    //obtain the query word(s)
    public void setup(Context context) throws IOException, InterruptedException {

	    Configuration conf = context.getConfiguration();
	    String temp[] = conf.get("keyword").split(" ");
	    for(int i = 0; i < temp.length; i++) keyword.add(temp[i]);
	    N = Double.parseDouble(conf.get("N"));
    }

    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {

	    String line = value.toString();
	    String info[] = line.split(";");
	    String temp[] = info[0].split("\\s+");
	    String term = temp[0];
        //ignore upper/lower case search
	    if(keyword.get(keyword.size()-1).equals("-case") == true) {
	        for(int j = 0; j < keyword.size()-1; j++) {
	            if(keyword.get(j).length() == 0) continue;
		        if(term.toLowerCase().equals(keyword.get(j).toLowerCase()) == true) {
	                context.getCounter(TOP10_COUNTER.FOUND).increment(1);
		            df = Double.parseDouble(temp[1]);
	                for(int k = 1; k < info.length; k++) {
			            String doc[] = info[k].split("\\s+");
			            int docID = Integer.valueOf(doc[0]);
			            Double tf = Double.parseDouble(doc[1]);
			            Double weight = tf * Math.log10(N/df);
			            context.write(new QueryKeyPair(new DoubleWritable(weight), docID), new QueryValuePair(docID, new Text(doc[2])));
		            }
		        }
	        }
	    }
        //case-sensitive search
	    else {
	        for(int j = 0; j < keyword.size(); j++) {
		        if(keyword.get(j).length() == 0) continue;
		        if(term.equals(keyword.get(j)) == true) {
		            context.getCounter(TOP10_COUNTER.FOUND).increment(1);
		            df = Double.parseDouble(temp[1]);
		            for(int k = 1; k < info.length; k++) {
			            String doc[] = info[k].split("\\s");
			            int docID = Integer.valueOf(doc[0]);
			            Double tf = Double.parseDouble(doc[1]);
			            Double weight = tf * Math.log10(N/df);
			            context.write(new QueryKeyPair(new DoubleWritable(weight), docID), new QueryValuePair(docID, new Text(doc[2])));
		            }
		        }
	        }
	    }
    }
}
