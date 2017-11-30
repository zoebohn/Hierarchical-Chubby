 package locks

import(
    "strings"
)

 func lock_array_to_string(lock_arr []Lock) string {
    var string_form []string
    for _, l := range lock_arr {
        string_form = append(string_form, string(l))
    }
    return strings.Join(string_form, ";")
 }

 func string_to_lock_array(s string) []Lock {
    string_arr := strings.Split(s, ";")
    var lock_arr []Lock
    for _, l := range string_arr {
        lock_arr = append(lock_arr, Lock(l))
    }
    return lock_arr
 }
