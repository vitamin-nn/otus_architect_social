  INSTALL PLUGIN rpl_semi_sync_master SONAME 'semisync_master.so';
  SET GLOBAL rpl_semi_sync_master_enabled = 1;
  SET GLOBAL rpl_semi_sync_master_wait_for_slave_count = 1;
  SHOW VARIABLES LIKE 'rpl_semi_sync%';