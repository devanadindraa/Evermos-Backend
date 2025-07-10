-- TABEL USER
CREATE TABLE
    user (
        id INT AUTO_INCREMENT PRIMARY KEY,
        nama VARCHAR(255),
        kata_sandi VARCHAR(255),
        notelp VARCHAR(255) UNIQUE,
        tanggal_lahir DATE,
        jenis_kelamin VARCHAR(255),
        tentang TEXT,
        pekerjaan VARCHAR(255),
        email VARCHAR(255) UNIQUE,
        id_provinsi VARCHAR(255),
        id_kota VARCHAR(255),
        isAdmin BOOLEAN,
        updated_at_date DATETIME,
        created_at_date DATETIME
    );

-- TABEL ALAMAT
CREATE TABLE
    alamat (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_user INT,
        judul_alamat VARCHAR(255),
        nama_penerima VARCHAR(255),
        no_telp VARCHAR(255),
        detail_alamat VARCHAR(255),
        updated_at_date DATETIME,
        created_at_date DATETIME,
        FOREIGN KEY (id_user) REFERENCES user (id)
    );

-- TABEL TOKO
CREATE TABLE
    toko (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_user INT,
        nama_toko VARCHAR(255),
        url_foto VARCHAR(255),
        updated_at_date DATETIME,
        created_at_date DATETIME,
        FOREIGN KEY (id_user) REFERENCES user (id)
    );

-- TABEL CATEGORY
CREATE TABLE
    category (
        id INT AUTO_INCREMENT PRIMARY KEY,
        nama_category VARCHAR(255),
        created_at_date DATETIME,
        updated_at_date DATETIME
    );

-- TABEL PRODUK
CREATE TABLE
    produk (
        id INT AUTO_INCREMENT PRIMARY KEY,
        nama_produk VARCHAR(255),
        slug VARCHAR(255),
        harga_reseller VARCHAR(255),
        harga_konsumen VARCHAR(255),
        stok INT,
        deskripsi TEXT,
        created_at_date DATETIME,
        updated_at_date DATETIME,
        id_toko INT,
        id_category INT,
        FOREIGN KEY (id_toko) REFERENCES toko (id),
        FOREIGN KEY (id_category) REFERENCES category (id)
    );

-- TABEL FOTO PRODUK
CREATE TABLE
    foto_produk (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_produk INT,
        url VARCHAR(255),
        updated_at_date DATETIME,
        created_at_date DATETIME,
        FOREIGN KEY (id_produk) REFERENCES produk (id)
    );

-- TABEL LOG PRODUK
CREATE TABLE
    log_produk (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_produk INT,
        nama_produk VARCHAR(255),
        slug VARCHAR(255),
        harga_reseller VARCHAR(255),
        harga_konsumen VARCHAR(255),
        deskripsi TEXT,
        created_at_date DATETIME,
        updated_at_date DATETIME,
        id_toko INT,
        id_category INT,
        FOREIGN KEY (id_toko) REFERENCES toko (id),
        FOREIGN KEY (id_category) REFERENCES category (id)
    );

-- TABEL TRANSAKSI
CREATE TABLE
    trx (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_user INT,
        alamat_pengiriman INT,
        harga_total INT,
        kode_invoice VARCHAR(255),
        method_bayar VARCHAR(255),
        updated_at_date DATETIME,
        created_at_date DATETIME,
        FOREIGN KEY (id_user) REFERENCES user (id),
        FOREIGN KEY (alamat_pengiriman) REFERENCES alamat (id)
    );

-- TABEL DETAIL TRANSAKSI
CREATE TABLE
    detail_trx (
        id INT AUTO_INCREMENT PRIMARY KEY,
        id_trx INT,
        id_log_produk INT,
        id_toko INT,
        kuantitas INT,
        harga_total INT,
        update_at_date DATETIME,
        created_at_date DATETIME,
        FOREIGN KEY (id_trx) REFERENCES trx (id),
        FOREIGN KEY (id_log_produk) REFERENCES log_produk (id),
        FOREIGN KEY (id_toko) REFERENCES toko (id)
    );

-- TABEL INVALID TOKEN
CREATE TABLE invalid_token (
    token VARCHAR(255) PRIMARY KEY,
    expires TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
