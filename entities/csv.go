package entities

// Result01 is csv
type Result01 struct {
	JobName              string `csv:"job_name"`
	JobID                string `csv:"job_ID"`
	TaskID               string `csv:"task_ID"`
	Type                 string `csv:"type"`
	StartTime            string `csv:"start_time"`
	FinishTime           string `csv:"finish_time"`
	ElapsedTime          string `csv:"elapsed_time"`
	Hostname             string `csv:"host"`
	FileBytesRead        string `csv:"FILE_BYTES_READ"`
	FileBytesWritten     string `csv:"FILE_BYTES_WRITTEN"`
	FileReadOps          string `csv:"FILE_READ_OPS"`
	FileLargeReadOps     string `csv:"FILE_LARGE_READ_OPS"`
	FileWriteOps         string `csv:"FILE_WRITE_OPS"`
	HDFSBytesRead        string `csv:"HDFS_BYTES_READ"`
	HDFSBytesWritten     string `csv:"HDFS_BYTES_WRITTEN"`
	HDFSReadOps          string `csv:"HDFS_READ_OPS"`
	HDFSLargeReadOps     string `csv:"HDFS_LARGE_READ_OPS"`
	HDFSWriteOps         string `csv:"HDFS_WRITE_OPS"`
	CombineInputRecords  string `csv:"COMBINE_INPUT_RECORDS"`
	CombineOutputRecords string `csv:"COMBINE_OUTPUT_RECORDS"`
	ReduceInputGroups    string `csv:"REDUCE_INPUT_GROUPS"`
	ReduceShuffleBytes   string `csv:"REDUCE_SHUFFLE_BYTES"`
	ReduceInputRecords   string `csv:"REDUCE_INPUT_RECORDS"`
	ReduceOutputRecords  string `csv:"REDUCE_OUTPUT_RECORDS"`
	SpilledRecords       string `csv:"SPILLED_RECORDS"`
	ShuffledMaps         string `csv:"SHUFFLED_MAPS"`
	FailedShuffle        string `csv:"FAILED_SHUFFLE"`
	MergedMapOutputs     string `csv:"MERGED_MAP_OUTPUTS"`
	GCTimeMillis         string `csv:"GC_TIME_MILLIS"`
	CPUMilliseconds      string `csv:"CPU_MILLISECONDS"`
	PhysicalMemoryBytes  string `csv:"PHYSICAL_MEMORY_BYTES"`
	VirtualMemoryBytes   string `csv:"VIRTUAL_MEMORY_BYTES"`
	CommitedHeapBytes    string `csv:"COMMITTED_HEAP_BYTES"`
	BadID                string `csv:"BAD_ID"`
	Connection           string `csv:"CONNECTION"`
	IOError              string `csv:"IO_ERROR"`
	WrongLength          string `csv:"WRONG_LENGTH"`
	WrongMap             string `csv:"WRONG_MAP"`
	WrongReduce          string `csv:"WRONG_REDUCE"`
	BytesWritten         string `csv:"BYTES_WRITTEN"`
}
