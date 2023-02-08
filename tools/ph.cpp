#include <iostream>
#include <Windows.h>

using namespace std;

int main()
{
    // HANDLE hProcess;
    // HANDLE hThread;
    // DWORD dwProcessId;
    // DWORD dwThreadId;
    STARTUPINFOA si;
    PROCESS_INFORMATION pi;
    ZeroMemory(&si, sizeof(si));
    ZeroMemory(&pi, sizeof(pi));
    bool bCreateProcess;

    bCreateProcess = CreateProcess("C:/Users/X-alien/AppData/Local/Programs/Ankama Launcher/Ankama Launcher.exe", NULL, NULL, NULL, FALSE, CREATE_SUSPENDED, NULL, NULL, &si, &pi);

    if (bCreateProcess == FALSE)
    {
        cout << "process creation failed ! " << GetLastError() << endl;
        return -1;
    }

    cout << "process creation succeded ! " << GetProcessId(pi.hProcess) << endl;
    cout << "process creation succeded ! " << GetThreadId(pi.hThread) << endl;

    ResumeThread(pi.hThread);
}