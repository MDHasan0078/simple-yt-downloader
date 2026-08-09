; Inno Setup script: Simple YT Downloader (Flutter app + engine) Windows installer.
; BUNDLE points to the Flutter release folder (app/build/windows/x64/runner/Release),
; which contains simple_yt_downloader.exe, data/, lib/, and engine/.
; Build from CI:
;   ISCC.exe app/build/windows/installer.iss /DVERSION=2.1.0 /DBUNDLE=<release> /DOUTPUT=<out>
;
; Output: <OUTPUT>/simple-yt-downloader-<VERSION>-Setup.exe

#ifndef VERSION
  #define VERSION "2.1.0"
#endif
#ifndef BUNDLE
  #define BUNDLE "..\build\windows\x64\runner\Release"
#endif
#ifndef OUTPUT
  #define OUTPUT "..\build"
#endif
#ifndef ICON
  #define ICON "runner\resources\app_icon.ico"
#endif

#define MyAppName "Simple YT Downloader"
#define MyAppPublisher "MD. Mahmudul Hasan"
#define MyAppVersion VERSION
#define MyAppExeName "simple_yt_downloader.exe"

[Setup]
AppId={{B8A9F1D2-4C3E-4A5B-9F0E-2C1D3E4F5A6C}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\SimpleYTDownloader
DefaultGroupName=Simple YT Downloader
DisableProgramGroupPage=yes
OutputDir={#OUTPUT}
OutputBaseFilename=simple-yt-downloader-{#MyAppVersion}-Setup
SourceDir={#BUNDLE}
Compression=lzma2
SolidCompression=yes
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
WizardStyle=modern
UninstallDisplayName=Simple YT Downloader
SetupIconFile={#ICON}
UninstallDisplayIcon={app}\{#MyAppExeName}

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Flutter app (exe, DLLs, data/) — everything except the engine folder
Source: "*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs; Excludes: "engine"
; Engine + bundled binaries
Source: "engine\*"; DestDir: "{app}\engine"; Flags: ignoreversion recursesubdirs createallsubdirs

[Icons]
Name: "{autoprograms}\Simple YT Downloader"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\Simple YT Downloader"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,Simple YT Downloader}"; Flags: nowait postinstall skipifsilent
