#### 通用建筑分类存储结构

建筑包含基本的坐标信息

建筑类别的区分用TYPE，TYPE值为32位Hash的字符串，即:

naming.FNV1a("某某建筑")

单个建筑的私有数据用一个KV表存储:  

map[string]string
