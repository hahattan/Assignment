import org.apache.spark._
import org.apache.hadoop.fs._
import scala.collection._

object ReachabilityQuery {
    def main(args: Array[String]) {
        val filePath = args(0)
        val name = args(1)
        val iters = args(2).toInt

        val outputPath = "Lab6/reachability"

        val conf = new SparkConf().setAppName("Reachability Query Application")
        val sc = new SparkContext(conf)

        // Cleanup output dir
        val hadoopConf = sc.hadoopConfiguration
        var hdfs = FileSystem.get(hadoopConf)
        try { hdfs.delete(new Path(outputPath), true) } catch { case _ : Throwable => { } }

        // Read input file
        val lines = sc.textFile(filePath, sc.defaultParallelism)

        // Step 1: initial friend list, map each line in file to pairs
        // ex: "Tom,Alice" => ("Tom" -> "Alice"), ("Alice" -> "Tom")
        // hint: use map() and union()
        //       Or, you can use flatMap() to do it at once
        val friend = lines.flatMap (line => {
            val names = line.split(",")
	    val t1 = (names(0), names(1))
	    val t2 = (names(1), names(0))
	    List(t1, t2)
        })

        // Step 2: make an RDD from 'name'
        // hint: use sc.parallelize()
        // optional: you can pass (sc.defaultParallelism * 3) as # of slices for better performance
        var res = sc.parallelize(Seq(name), 3*sc.defaultParallelism)

        // Step 3: use initial result and friend list to find all friends within 'iters' hops
        // hint: use join(), union() and distinct()
        // optional: resize numPartitions to (sc.defaultParallelism * 3) at the end of each iteration
        for (i <- 1 to iters) {
	    var temp = res.map(x => (x, x))
	    val joined = temp.join(friend)
	    val temp2 = joined.map(x => x._2._2)
            res = res.union(temp2)
            res = res.distinct()
        }

        // Action branch! Add cache() to avoid re-computation
        res = res.cache

        // Output: count() and saveAsTextFile()
        println("Count := " + res.count)
        res.sortBy(_.toString).saveAsTextFile(outputPath)

        sc.stop
    }
}
