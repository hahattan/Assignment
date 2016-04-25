package invertedIndex;

import org.apache.hadoop.io.Text;
import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.WritableComparator;

public class InvertedIndexGroupComparator extends WritableComparator {

  public InvertedIndexGroupComparator() {
    super(InvertedIndexKeyPair.class, true);
  }


  public int compare(WritableComparable o1, WritableComparable o2) {
    InvertedIndexKeyPair key1 = (InvertedIndexKeyPair) o1;
    InvertedIndexKeyPair key2 = (InvertedIndexKeyPair) o2;

    return key1.compareTo(key2);
  }
}
