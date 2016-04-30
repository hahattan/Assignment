package query;

import java.net.URI;
import java.io.*;
import java.util.*;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.*;
import org.apache.hadoop.util.*;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapred.Counters;

enum TOP10_COUNTER {
    NUMBER,
    RANK,
    FOUND
};

public class Query {

    public static void main(String[] args) throws Exception {

	System.out.println("Enter the word you want to look up : ");
	Scanner input = new Scanner(System.in);
	String keyword = input.nextLine();
	input.close();

	Configuration conf = new Configuration();
	conf.set("keyword", keyword);
	conf.set("path", args[1]);

	// get the number of input files 
	try {
	    FileSystem fs = FileSystem.get(conf);
	    Path path = new Path("InvertedIndex/id.log");
	    BufferedReader br = new BufferedReader(new InputStreamReader(fs.open(path)));
	    String line;
	    line = br.readLine();
	    conf.set("N", line);
	}catch(Exception e) {

	}


	Job job = Job.getInstance(conf, "Retrieval");
	job.setJarByClass(Query.class);

	// set the class of each stage in mapreduce
	job.setMapperClass(QueryMapper.class);
	job.setSortComparatorClass(QueryKeyComparator.class);
	job.setGroupingComparatorClass(QueryGroupComparator.class);
	job.setReducerClass(QueryReducer.class);

	// set the output class of Mapper and Reducer
	job.setMapOutputKeyClass(QueryKeyPair.class);
	job.setMapOutputValueClass(QueryValuePair.class);
	job.setOutputKeyClass(Text.class);
	job.setOutputValueClass(Text.class);

	// set the number of reducer
	job.setNumReduceTasks(1);

	// add input/output path
	FileInputFormat.addInputPath(job, new Path(args[0]));
	FileOutputFormat.setOutputPath(job, new Path(args[2]));

	System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
