set style line 1 lt 1 pt 0;
set title "Books needed to learn English";
set label "25,000 words: 99.99" at 25000,99.99 point lt 0 pt 9 offset -9,1;
set label "30,000 words: 122.2" at 30000,122.2 point lt 0 pt 9 offset -9,1;
set label "35,000 words: 143.6" at 35000,143.6 point lt 0 pt 9 offset -9,1;
set xlabel "Vocabulary size";
set ylabel "Number of books (simulated)";
set term png
set output '| display png:-'
plot "data.txt" using 1:3 with line ls 1 notitle
