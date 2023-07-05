end=0
TIMEFORMAT=%R
while [ 100000 -ge $end ]
do
#duracion=$( TIMEFORMAT="%R"; { time (ls 2>/dev/null); } ) 

	start=$(date +%s%3N)
        curl -s "localhost:7043/uploader/state?hash=MXxmYWtlVXNlcm5hbWU="
	elapsed=$(expr $(date +%s%3N) - $start)
	now=$(date)
	echo "duration ${elapsed} time ${now}"
	let end=$end+1
	sleep 2
done
