/*
* @Author: Ins
* @Date:   2018-10-10 09:54:12
* @Last Modified by:   Ins
* @Last Modified time: 2018-10-16 16:29:06
*/
package main
import "C"
import (
    "fmt"
    "unsafe"
    "github.com/ceph/go-ceph/rados"
)


func newConn(cluster_name, user_name, conf_file string) (*rados.Conn, error) {
    conn, err := rados.NewConnWithClusterAndUser(cluster_name, user_name)//"ceph","client.objstore"
    if err != nil {
        return nil, err
    }

    err = conn.ReadConfigFile(conf_file)//"/etc/ceph/ceph.conf"
    if err != nil {
        return nil, err
    }

    err = conn.Connect()
    if err != nil {
        return nil, err
    }

    return conn, nil
}

func ReadObjectToBytes(ioctx *rados.IOContext, oname string, block_size int, offset uint64) (int, []byte, error) {
    bytesOut := make([]byte, block_size)
    ret, err := ioctx.Read(oname, bytesOut, offset)
    if err != nil {
        return -1, bytesOut, err
    }
    bytesOut = bytesOut[:ret]
    return ret, bytesOut, err
}

func ObjectListFunc(oid string) {
    fmt.Println(oid)
}

//export ListObj
func ListObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char) *C.char{
    cluster_name, user_name, conf_file, pool_name := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext" + err.Error())
    }
    defer ioctx.Destroy()

    // list the objects in pool just printed in terminal
    ioctx.ListObjects(ObjectListFunc)
    return C.CString("list the objects above in object:" + pool_name)
}

//export FromObj
func FromObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char,block_size int, c_oname *C.char, offset uint64) *C.char{
    if block_size > 204800000 {
        return C.CString("the block_size cannot be greater than 204800000")
    }
    cluster_name, user_name, conf_file, pool_name, oname := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oname)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext" + err.Error())
    }
    defer ioctx.Destroy()

    // read the data and write to the file

    ret, bytesOut, err := ReadObjectToBytes(ioctx, oname, block_size, offset)
    if ret == -1 {
        return C.CString("error when read the object to bytes:" + err.Error())
    }
    
    return C.CString(*(*string)(unsafe.Pointer(&bytesOut)))

}

//export WriteToObj
func WriteToObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oname *C.char, c_bytesIn *C.char, offset uint64) *C.char{
    cluster_name, user_name, conf_file, pool_name, oname, bytesIn := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oname), C.GoString(c_bytesIn)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext:" + err.Error())
    }
    defer ioctx.Destroy()

    // write data to object
    err = ioctx.Write(oname, []byte(bytesIn), offset)
    if err != nil {
        return C.CString("error when write to object:" + err.Error())
    }

    return C.CString("successfully writed to object：" + oname)
}

//export AppendToObj
func AppendToObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oname *C.char, c_bytesIn *C.char) *C.char{
    cluster_name, user_name, conf_file, pool_name, oname, bytesIn := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oname), C.GoString(c_bytesIn)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext:" + err.Error())
    }
    defer ioctx.Destroy()

    // write data to object
    err = ioctx.Append(oname, []byte(bytesIn))
    if err != nil {
        return C.CString("error when append to object:" + err.Error())
    }

    return C.CString("successfully append to object：" + oname)
}

//export DelObj
func DelObj(c_cluster_name *C.char, c_user_name *C.char, c_conf_file *C.char, c_pool_name *C.char, c_oname *C.char) *C.char{
    cluster_name, user_name, conf_file, pool_name, oname := C.GoString(c_cluster_name), C.GoString(c_user_name), C.GoString(c_conf_file), C.GoString(c_pool_name), C.GoString(c_oname)
    conn, err := newConn(cluster_name, user_name, conf_file)
    if err != nil {
        return C.CString("error when invoke a new connection:" + err.Error())
    }
    defer conn.Shutdown()

    // open a pool handle
    ioctx, err := conn.OpenIOContext(pool_name)
    if err != nil {
        return C.CString("error when openIOContext" + err.Error())
    }
    defer ioctx.Destroy()

    // delete a object 
    err = ioctx.Delete(oname)
    if err != nil {
        return C.CString("error when delete the object:" + err.Error())
    }
    return C.CString("successfully delete the object:" + oname)
}
func main() {
    
}