{ config, pkgs, lib, ... }:
with lib;
let cfg = config.services.vpub-plus-plus;
in {

  options.services.vpub-plus-plus = {
    enable = mkEnableOption "vpub-plus-plus service";

    title = mkOption {
      type = types.str;
      default = null;
      example = "My Forum";
      description = "Title of the forum";
    };

    address = mkOption {
      type = types.str;
      default = "127.0.0.1";
      example = "0.0.0.0";
      description = "Address to listen on";
    };

    port = mkOption {
      type = types.str;
      default = "1234";
      example = "2345";
      description = "Port to listen on";
    };

    DBUri = mkOption {
      type = types.str;
      default = "postgres:///vpub?host=/run/postgresql&sslmode=disable";
      description = "Postgres connection URI";
    };

    envFile = mkOption {
      type = types.str;
      default = null;
      example = "/var/secrets/vpub-plus-plus/envfile";
      description = ''
        Additional environment file to pass to the service, containing:
        32 bytes long session key (SESSION_KEY) and 32 bytes long CSRF key (CSRF_KEY)
      '';
    };
  };

  config = mkIf cfg.enable {

    # User and group
    users.users.vpub = {
      isSystemUser = true;
      description = "vpub-plus-plus user";
      extraGroups = [ "vpub" ];
      group = "vpub";
    };

    users.groups.vpub-plus-plus.name = "vpub";

    # Service
    systemd.services.vpub-plus-plus = {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ];
      description = "vpub-plus-plus";
      environment = {
          # Postgresql connection URL
          DATABASE_URL = cfg.DBUri;
          # Your forum name
          TITLE = cfg.title;
          # What port is going to be used by a vpub HTTP server
          PORT = cfg.port;
          # Address to listen on
          ADDRESS= cfg.address;
        };
      serviceConfig = {

        EnvironmentFile = [ cfg.envFile ];

        User = "vpub";
        ExecStart = "${pkgs.vpub-plus-plus}/bin/vpub";
        Restart = "on-failure";
        RestartSec = "5s";
      };
    };

    # TODO Postgres config
    services.postgresql = {
      enable = true;

      ensureUsers = [{
        name = "vpub";
        ensurePermissions = {
          "DATABASE vpub" = "ALL PRIVILEGES";
        };
      }];

      ensureDatabases = [ "vpub" ];
    };


    # TODO nginx config

    # services.nginx = {
    #   enable = true;
    #   virtualHosts."server" = {
    #     root = pkgs.vpub-plus-plus;
    #   };
    # };
  };
}
