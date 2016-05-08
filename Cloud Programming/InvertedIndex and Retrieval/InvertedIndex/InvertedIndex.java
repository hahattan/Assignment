package invertedIndex;

import java.net.URI;
import java.io.*;
import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.*;
import org.apache.hadoop.util.*;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class InvertedIndex {

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();

        Job job = Job.getInstance(conf, "InvertedIndex");
        job.setJarByClass(InvertedIndex.class);

        // set the class of each stage in mapreduce
        job.setMapperClass(InvertedIndexMapper.class);
        job.setCombinerClass(InvertedIndexCombiner.class);
        //job.setPartitionerClass(InvertedIndexPartitioner.class);
        job.setSortComparatorClass(InvertedIndexKeyComparator.class);
        job.setGroupingComparatorClass(InvertedIndexGroupComparator.class);
        job.setReducerClass(InvertedIndexReducer.class);

        // set the output class of Mapper and Reducer
        job.setMapOutputKeyClass(InvertedIndexKeyPair.class);
        job.setMapOutputValueClass(InvertedIndexValuePair.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class);

        // set the number of reducer
        job.setNumReduceTasks(1);

        // add input/output path
        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        // write all input file names to a file to get docID
        try {
            FileSystem fs = FileSystem.get(URI.create(args[1]), conf);
            Path file = new Path("id.log");
            FileStatus[] status = fs.listStatus(new Path(args[0]));
            BufferedWriter bw = new BufferedWriter(new OutputStreamWriter(fs.create(file, true)));
            bw.write(String.valueOf(status.length));
            bw.write("\n");
            System.out.println(status.length);
            for(int i = 0; i < status.length; i++) {
	            bw.write(status[i].getPath().getName());
	            bw.write("\n");
            }
            bw.close();
            job.addCacheFile(file.toUri());
        }catch(Exception e) {

        }
        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
