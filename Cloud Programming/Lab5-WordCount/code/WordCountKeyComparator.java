package wordcount;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class WordCountKeyComparator extends WritableComparator {

  public WordCountKeyComparator() {
    super(Text.class, true);
  }


  public int compare(WritableComparable o1, WritableComparable o2) {
    Text key1 = (Text) o1;
    Text key2 = (Text) o2;
    int value = 0;
    int a1, a2, f1 = 0, f2 = 0;
    // TODO Order by A -> a -> B -> b ....
    if(key1.charAt(0) >= 97) {
	a1 = key1.charAt(0)-32;
	f1 = 1;
    }
    else a1 = key1.charAt(0);
    
    if(key2.charAt(0) >= 97) {
	a2 = key2.charAt(0)-32;
	f2 = 1;
    }
    else a2 = key2.charAt(0);

    if(a1 < a2) value = -1;
    else if(a1 > a2) value = 1;
    else if(a1 == a2) {
	if(f1 == f2) value = 0;
	else if(f1 == 0 && f2 == 1) value = -1;
	else if(f1 == 1 && f2 == 0) value = 1;
    }
    return value;
  }
}  
