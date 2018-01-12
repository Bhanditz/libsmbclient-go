#include "libsmbclient.h"
#include "_cgo_export.h"



int my_errno() {
    return errno;
}


/***********************************************************************************************************************
    CALLABLE FUNCTIONS FOR FILES
***********************************************************************************************************************/


/*
 SMB client AUTHENTICATE
*/
void my_smbc_auth_callback(SMBCCTX *ctx,
	       const char *server_name, const char *share_name,
	       char *workgroup_out, int workgroup_len,
	       char *username_out, int username_len,
	       char *password_out, int password_len) {

   void *go_fn = smbc_getOptionUserData(ctx);
   GoAuthCallbackHelper(
   go_fn,
   (char*)server_name,
   (char*)share_name,
   workgroup_out,
   workgroup_len,
   username_out,
   username_len,
   password_out,
   password_len);
}

/*
 SMB client AUTH INIT
*/
void my_smbc_init_auth_callback(SMBCCTX *ctx, void *go_fn){
   smbc_setOptionUserData(ctx, go_fn);
   smbc_setFunctionAuthDataWithContext(ctx, my_smbc_auth_callback);
}




///*
// SMB client SET NOTIFY FUNCTION
//*/
//int my_smbc_notify_callback(
//    const struct smbc_notify_callback_action *actions,
//    size_t num_actions,
//    void *private_data
//) {
//    //GoChangesNotifyCallback(actions, private_data);
//}
//
//
///*
// SMB client CHANGES INIT
//*/
//void my_smbc_init_notify_callback(SMBCCTX *ctx, void *go_fn) {
//    smbc_setFunctionNotify(ctx, my_smbc_notify_callback);
//}


/**********************************************************************************************************************
    CALLABLE FUNCTIONS FOR DIRECTORIES
***********************************************************************************************************************/

/*
 SMB client CREATE
*/
int my_smbc_mkdir(SMBCCTX *c, const char *dirname, mode_t mode) {
    smbc_mkdir_fn fn = smbc_getFunctionMkdir(c);
    return fn(c, dirname, mode);
}


/*
 SMB client OPEN
*/
SMBCFILE* my_smbc_opendir(SMBCCTX *c, const char *fname) {
  smbc_opendir_fn fn = smbc_getFunctionOpendir(c);
  return fn(c, fname);
}


/*
 SMB client READ
*/
struct smbc_dirent* my_smbc_readdir(SMBCCTX *c, SMBCFILE *dir) {
  smbc_readdir_fn fn = smbc_getFunctionReaddir(c);
  return fn(c, dir);
}


/*
 SMB client CLOSE
*/
int my_smbc_closedir(SMBCCTX *c, SMBCFILE *dir) {
  smbc_closedir_fn fn = smbc_getFunctionClosedir(c);
  return fn(c, dir);
}





/***********************************************************************************************************************
    CALLABLE FUNCTIONS FOR FILES
***********************************************************************************************************************/

/*
 SMBC client OPEN-FILE
*/
SMBCFILE* my_smbc_open(SMBCCTX *c, const char *fname, int flags, mode_t mode) {
  smbc_open_fn fn = smbc_getFunctionOpen(c);
  return fn(c, fname, flags, mode);
}


/*
 SMB client CREATE-FILE
*/
SMBCFILE* my_smbc_creat(SMBCCTX *c, const char *fname, mode_t mode) {
    smbc_creat_fn fn = smbc_getFunctionCreat(c);
    return fn(c, fname, mode);
}


/*
 SMB client READ
*/
ssize_t my_smbc_read(SMBCCTX *c, SMBCFILE *file, void *buf, size_t count) {
  smbc_read_fn fn = smbc_getFunctionRead(c);
  return fn(c, file, buf, count);
}


/*
 SMB client WRITE
*/
ssize_t my_smbc_write(SMBCCTX *c, SMBCFILE *file, void *buf, size_t count) {
  smbc_write_fn fn = smbc_getFunctionWrite(c);
  return fn(c, file, buf, count);
}


/*
 SMB client UNLINK
*/
ssize_t my_smbc_unlink(SMBCCTX *c, const char *url) {
  smbc_unlink_fn fn = smbc_getFunctionUnlink(c);
  return fn(c, url);
}


/*
 SMB client RENAME
*/
int my_smbc_rename(SMBCCTX *c, const char *ourl, SMBCCTX *nc, const char *nurl) {
  smbc_rename_fn fn = smbc_getFunctionRename(c);
  return fn(c, ourl, nc, nurl);
}


/*
 SMB client LSEEK
*/
off_t my_smbc_lseek(SMBCCTX *c, SMBCFILE * file, off_t offset, int whence) {
  smbc_lseek_fn fn = smbc_getFunctionLseek(c);
  return fn(c, file, offset, whence);
}


/*
 SMB client STAT
*/
int my_smbc_stat(SMBCCTX *c, const char *url, struct stat *st) {
    smbc_stat_fn fn = smbc_getFunctionStat(c);
    return fn(c, url, st);
}


/*
 SMB client CLOSE
*/
int my_smbc_close(SMBCCTX *c, SMBCFILE *f) {
  smbc_close_fn fn = smbc_getFunctionClose(c);
  fn(c, f);
}