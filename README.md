# plugin-postgresql
 

**利用sqleyes引擎实现的postgresql 数据库的监控插件**



> 抓包截取项目中的postgresql数据库请求并解析成相应的SQL语句便于调试。 不需要修改代码，直接嗅探项目中的数据请求。



## Support List:

*Client*

- [x] p ->Authentication message
- [x] Q -> Simple query

*Server*

- [x] R -> Authentication request
- [x] K -> Backend key data
- [x] S -> Parameter status
- [x] C -> Command completion
- [x] T -> Row description
- [x] D -> Data row
- [x] E -> Error
- [x] N -> Notice
- [x] Z -> Ready for query

## TODO：

*Client*

- [ ] P -> Parse
- [ ] B -> Bind
- [ ] E -> Execute
- [ ] D -> Describe
- [ ] C -> Close
- [ ]  H -> Flush
- [ ]  S -> Sync
- [ ]  F -> Function call
- [ ]  d -> Copy data
- [ ]  c -> Copy completion
- [ ]  f -> Copy failure
- [ ]  X -> Termination  

*Server*

- [ ] 1 -> Parse completion
- [ ] 2 -> Bind completion
- [ ] 3 -> Close completion
- [ ] t -> Parameter description
- [ ] I -> Empty query
- [ ] n -> No data
- [ ] s -> Portal suspended
- [ ] A -> Notification
- [ ] V -> Function call response
- [ ] G -> CopyIn response
- [ ] H -> CopyOut response
- [ ] d -> Copy data
- [ ] c -> Copy completion
- [ ] v -> Negotiate protocol version 