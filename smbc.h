#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <libsmbclient.h>

typedef struct stat c_stat;



int my_errno();


/***********************************************************************************************************************
    CALLABLE FUNCTIONS FOR FILES
***********************************************************************************************************************/

/*
 the auth callback related stuff, its a bit complicated
 SMB client AUTHENTICATE HELPER
*/
extern void GoAuthCallbackHelper(
    void *o,
    char *server_name,
    char *share_name,
    char *domain_out, int domain_len,
    char *username_out, int username_len,
    char *password_out, int password_len
);

/*
 SMB client AUTHENTICATE
*/
void my_smbc_auth_callback(SMBCCTX *o,
           const char *server_name, const char *share_name,
           char *domain_out, int domain_len,
           char *username_out, int username_len,
           char *password_out, int password_len);

/*
 SMB client AUTH INIT
*/
void my_smbc_init_auth_callback(SMBCCTX *ctx, void *go_fn);




/*
 the changes callback
 SMB client CHANGES NOTIFY
*/
//extern void GoChangesNotifyCallback(const struct smbc_notify_callback_action* actions, void *private_data);

///*
// SMB client AUTHENTICATE
//*/
//void my_smbc_notify_callback(
//    const struct smbc_notify_callback_action *actions,
//    size_t num_actions,
//    void *private_data);
//
///*
// SMB client CHANGES INIT
//*/
//void my_smbc_init_notify_callback(SMBCCTX *ctx, void *go_fn);





/**********************************************************************************************************************
    CALLABLE FUNCTIONS FOR DIRECTORIES
***********************************************************************************************************************/

/*
 SMB client CREATE
*/
int my_smbc_mkdir(SMBCCTX *c, const char *dirname, mode_t mode);


/*
 SMB client OPEN
*/
SMBCFILE* my_smbc_opendir(SMBCCTX *c, const char *fname);


/*
 SMB client READ
*/
struct smbc_dirent* my_smbc_readdir(SMBCCTX *c, SMBCFILE *dir);


/*
 SMB client CLOSE
*/
int my_smbc_closedir(SMBCCTX *c, SMBCFILE *dir);





/***********************************************************************************************************************
    CALLABLE FUNCTIONS FOR FILES
***********************************************************************************************************************/

/*
 SMBC client OPEN-FILE
*/
SMBCFILE* my_smbc_open(SMBCCTX *c, const char *fname, int flags, mode_t mode);


/*
 SMB client CREATE-FILE
*/
SMBCFILE* my_smbc_creat(SMBCCTX *c, const char *fname, mode_t mode);


/*
 SMB client READ
*/
ssize_t my_smbc_read(SMBCCTX *c, SMBCFILE *file, void *buf, size_t count);


/*
 SMB client WRITE
*/
ssize_t my_smbc_write(SMBCCTX *c, SMBCFILE *file, void *buf, size_t count);


/*
 SMB client UNLINK
*/
ssize_t my_smbc_unlink(SMBCCTX *c, const char *url);


/*
 SMB client RENAME
*/
int my_smbc_rename(SMBCCTX *c, const char *ourl, SMBCCTX *nc, const char *nurl);


/*
 SMB client LSEEK
*/
off_t my_smbc_lseek(SMBCCTX *c, SMBCFILE * file, off_t offset, int whence);


/*
 SMB client STAT
*/
int my_smbc_stat(SMBCCTX *c, const char *url, struct stat *st);



/*
 SMB client CLOSE
*/
int my_smbc_close(SMBCCTX *c, SMBCFILE *f);