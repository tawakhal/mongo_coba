# mongo_coba



# Tool mongo_viewer
1. noSQL Booster


# Leasson Learn
1. MongoDB akan melakukan create collection secara otomatis
2. _id adalah index yang akan otomatis tercreate langsung
3. index yang tercreate otomatis _id itu adalah int32
4. ketika melakukan update, tapi tidak menggunakan _id, maka bisa dilakukan. 
    Tapi, _id tidak boleh berubah. akan error seperti : write exception: write errors: [Performing an update on the path '_id' would modify the immutable field '_id']
5. Ketika ingin melakukan update, tapi datanya itu ada banyak yang sama, dari selector yang kita ingingkan (where).
    maka, akan terkena error seperti : write exception: write errors: [Performing an update on the path '_id' would modify the immutable field '_id']
    hal ini karena idnya itu banyak yang berbeda dari yang ingin diupdate.
6. Ketika ingin melakukan delete, yang dimana selectornya (where) itu mempunyai nilai yang sama. contoh : 
    [{"_id" : 2, "name" : "namanya"},{"_id" : 1, "name" : "namanya"}] 
    jika, kita hapus dimana selector (where) itu adalah "name" : "namanya". Maka, yang akan terhapus adalah data yang lama yaitu "_id" : 1.
    Jadi berhati"Lah..