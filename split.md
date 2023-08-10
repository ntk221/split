# SYNOPSIS
     split -d [-l line_count] [-a suffix_length] [file [prefix]]
     split -d -b byte_count[K|k|M|m|G|g] [-a suffix_length] [file [prefix]]
     split -d -n chunk_count [-a suffix_length] [file [prefix]]
     split -d -p pattern [-a suffix_length] [file [prefix]]

# Description
The split utility reads the given file and breaks it up into files of 1000 lines each(if no options are specified), leaving the file unchanged, It file is a single dash ('-') or absent, split reads from the standard input.

The options are as follows:

-l line_count
	Create split files line_count lines in length(-l でsplitできる行数を指定できるっぽい)

-n chunk_count
	Split file into chunk_count smaller files. The first n - 1 files will be of size (size of file / chunk_count ) and the last file will contain the remaining bytes. (splitした行を詰め込むファイル数を指定できる。各ファイルの大きさは file / chunk_count になる)

-b byte_count[K|k|M|m|G|g]
	Create split files byte_count bytes in length. If k or K is appended to the number, the file is split into byte_count kilobyte pieces. If m or M is appended to the number, the file is split into byte_count megabyte pieces. If g or G is appended to the number, the file is split into byte_count gigabyte pieces.


