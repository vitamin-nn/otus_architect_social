    INSTALL PLUGIN rpl_semi_sync_slave SONAME 'semisync_slave.so';
    SET GLOBAL rpl_semi_sync_slave_enabled = 1;
    SHOW VARIABLES LIKE 'rpl_semi_sync%';