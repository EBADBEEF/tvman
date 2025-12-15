{
  config,
  pkgs,
  lib,
  ...
}:
let
  inherit (lib)
    mkIf
    mkEnableOption
    mkMerge
    mkOption
    optionals
    types
    ;

  cfg = config.ebadf.tvman;

  defaultEnv = [
    "TVMAN_URL=unix://${cfg.socket}"
    "MENUMAN_FONTPATH=${pkgs.dejavu_fonts}/share/fonts/truetype/DejaVuSans.ttf"
  ]
  ++ optionals cfg.menu [ "MENUMAN_URL=unix://${cfg.menuSocket}" ];

  tvmanServiceConfig = {
    systemd.services.tvman = {
      wantedBy = [
        "multi-user.target"
        "dev-ttyTV.device"
      ];
      requires = [
        "tvman.socket"
        "dev-ttyTV.device"
        "pulse8-inputattach.service"
      ];
      after = [
        "tvman.socket"
        "dev-ttyTV.device"
        "pulse8-inputattach.service"
        "lircd.socket"
      ];
      serviceConfig = {
        DeviceAllow = [
          "char-cec rw"
          "char-lirc rw"
          "char-ttyUSB rw"
        ];
        Environment = defaultEnv;
        ExecStart = ''${pkgs.ebadf.tvman}/bin/tvman -rc '${cfg.rc}' -rcprotos '${cfg.rcprotos}' -verbose server'';
        Restart = "always";
        RestartSec = "10s";
      };
    };

    systemd.sockets.tvman = {
      wantedBy = [ "sockets.target" ];
      socketConfig = {
        ListenStream = cfg.socket;
      };
    };

    services.udev.extraRules = ''
      SUBSYSTEM=="tty", ENV{ID_MODEL}=="Keyspan_USA-19H", SYMLINK+="ttyTV", TAG+="systemd"
      SUBSYSTEM=="tty", KERNEL=="ttyACM[0-9]*", ATTRS{idVendor}=="2548", ATTRS{idProduct}=="1002", SYMLINK+="ttyCEC", TAG+="systemd",
    '';

    systemd.services."pulse8-inputattach" = {
      wantedBy = [ "dev-ttyCEC.device" ];
      bindsTo = [ "dev-ttyCEC.device" ];
      after = [ "dev-ttyCEC.device" ];
      serviceConfig = {
        ExecStart = "${pkgs.linuxConsoleTools}/bin/inputattach --pulse8-cec /dev/ttyCEC";
      };
    };
  };

  menumanServiceConfig = {
    systemd.services.menuman = {
      wantedBy = [ "graphical.target" ];
      requires = [ "menuman.socket" ];
      after = [
        "menuman.socket"
        "graphical.target"
      ];
      serviceConfig = {
        #TODO: get vars from autologin?
        DynamicUser = true;
        Environment = [
          "DISPLAY=:0"
          "DBUS_SESSION_BUS_ADDRESS=/dev/null"
          "SDL_VIDEODRIVER=x11"
        ]
        ++ defaultEnv;
        ExecStart = ''${pkgs.ebadf.tvman}/bin/menuman'';
        Restart = "always";
        RestartSec = "10s";
      };
    };

    systemd.sockets.menuman = {
      wantedBy = [ "sockets.target" ];
      socketConfig = {
        ListenStream = cfg.menuSocket;
      };
    };
  };
in
{
  options.ebadf.tvman = {
    enable = mkEnableOption "Enable TV Manager";
    rc = mkOption {
      type = types.str;
      description = "rc device uevent matcher";
      default = "";
    };
    rcprotos = mkOption {
      type = types.str;
      description = "rc protos (\"proto1 proto2\")";
      default = "nec sony";
    };
    socket = mkOption {
      type = types.path;
      description = "socket path";
      default = "/run/tvman/socket";
    };
    menu = mkEnableOption "Enable TV Menu";
    menuSocket = mkOption {
      type = types.path;
      description = "socket path";
      default = "/run/menuman/socket";
    };
  };
  config = mkMerge [
    (mkIf cfg.enable tvmanServiceConfig)
    (mkIf cfg.menu menumanServiceConfig)
  ];
}
