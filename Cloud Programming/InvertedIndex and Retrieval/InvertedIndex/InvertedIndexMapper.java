package invertedIndex;

import java.io.*;
import java.util.*;
import java.net.URI;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.io.*;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.filecache.DistributedCache;
import org.apache.hadoop.mapreduce.lib.input.FileSplit;

public class InvertedIndexMapper extends Mapper<LongWritable, Text, InvertedIndexKeyPair, InvertedIndexValuePair> {

    private List<String> fileList = new ArrayList<String>();
    private InvertedIndexKeyPair term = new InvertedIndexKeyPair();
    private InvertedIndexValuePair info = new InvertedIndexValuePair();

    //read the file we wrote at begining
    public void setup(Context context) throws IOException, InterruptedException {

        URI[] local = context.getCacheFiles();
        Path file = new Path(local[0].getPath());

        FileSystem fs = FileSystem.get(context.getConfiguration());
        BufferedReader br = new BufferedReader(new InputStreamReader(fs.open(file)));
        String line;
        line = br.readLine();
        System.out.println(line);
        while(line != null) {
            fileList.add(line);
  	        line = br.readLine();
        }
    }

    public void map(LongWritable key, Text value, Context context) throws IOException, InterruptedException {
        //obtain current filename and its corresponding docID
        String filename = ((FileSplit)context.getInputSplit()).getPath().getName();
        int id = 0;
        for(int i = 0; i < fileList.size(); i++) {
            if(fileList.get(i).equals(filename)) {
                id = i;
                break;
            }
        }
        // we simply use StringTokenizer to split the words for us.
        StringTokenizer itr = new StringTokenizer(value.toString());
        while (itr.hasMoreTokens()) {
            String word = itr.nextToken();
            //split on non-alphabetic characters
            String[] words = word.split("\\P{Alpha}", -1);
            for(int i = 0; i < words.length; i++) {
	            if(words[i] != null && !words[i].isEmpty()) {
	                term = new InvertedIndexKeyPair(new Text(words[i]), id);
	                info = new InvertedIndexValuePair(id, new IntWritable(1), new Text(key.toString()));

	                // create <K, V> pair
	                context.write(term, info);
	            }
            }
        }
    }
}
