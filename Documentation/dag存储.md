
### 表设计(edge)
```
DROP TABLE
IF
EXISTS edge;
CREATE TABLE
IF
NOT EXISTS edge (
id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
entry_edge_id int NOT NULL,
direct_edge_id int NOT NULL,
exit_edge_id int NOT NULL,
start_vertex varchar(36) NOT NULL,
end_vertex varchar(36) NOT NULL,
hops int NOT NULL,
source VARCHAR(30) NOT NULL
) ENGINE = INNODB DEFAULT CHARSET = utf8mb4;
```

### 插入
```
delimiter $
CREATE PROCEDURE proc_add ( 
IN StartVertexId VARCHAR ( 36 ), 
IN EndVertexId VARCHAR ( 36 ), 
IN Source VARCHAR ( 150 ) ) 
BEGIN
	IF
		EXISTS ( SELECT id FROM edge WHERE start_vertex = StartVertexId AND end_vertex = EndVertexId AND hops = 0 ) THEN
		 SELECT	'已存在';
     rollback;
	END IF;
	
	IF
		EXISTS ( SELECT id FROM edge WHERE start_vertex = EndVertexId AND end_vertex = StartVertexId ) THEN
		 SELECT	'存在环';		 
		 rollback;
	END IF;	

#插入新边
insert into edge (start_vertex, end_vertex, hops, source) 
values(StartVertexId, EndVertexId, 0 ,Source);

set @id = last_insert_id();

update edge 
  set entry_edge_id = @id,
	  exit_edge_id = @id,
	  direct_edge_id = @id
		where id = @id;

# step 1: A's incoming edges to B
insert into  edge (
  entry_edge_id,
  direct_edge_id,
	exit_edge_id,
	start_vertex,
	end_vertex,
	hops,
	source)
  select 
	 id,
	 @id,
	 @id,
	 start_vertex,
	 EndVertexId,
	 hops + 1,
	 Source
  from
	 edge
  where end_vertex = StartVertexId;

# step 2: A to B's outgoing edges
insert into  edge (
  entry_edge_id,  #作为这条隐含边的创建原因的起始顶点的传入边的 ID ；直接边包含与 Id 列相同的值
  direct_edge_id, #	导致创建此隐含边的直接边的 ID ；直接边包含与 Id 列相同的值
	exit_edge_id,   #作为这条隐含边的创建原因的结束顶点的出边的 ID ；直接边包含与 Id 列相同的值
	start_vertex,
	end_vertex,
	hops,
	source)
  select 
	 @id,
	 @id,
	 id,
	 StartVertexId,
	 end_vertex,
	 hops + 1,
	 Source
  from
	 edge
  where start_vertex = EndVertexId;

# A’s incoming edges to end vertex of B's outgoing edges
insert into  edge (
  entry_edge_id,  #作为这条隐含边的创建原因的起始顶点的传入边的 ID ；直接边包含与 Id 列相同的值
  direct_edge_id, #	导致创建此隐含边的直接边的 ID ；直接边包含与 Id 列相同的值
	exit_edge_id,   #作为这条隐含边的创建原因的结束顶点的出边的 ID ；直接边包含与 Id 列相同的值
	start_vertex,
	end_vertex,
	hops,
	source)
  select 
	 A.id,
	 @id,
	 B.id,
	 A.start_vertex,
	 B.end_vertex,
	 A.hops + B.hops + 1,
	 Source
  from edge A
  cross join
	 edge B
	where
		A.end_vertex = StartVertexId
    and B.start_vertex = EndVertexId;
END $
delimiter;
```

### 测试数据
```bigquery
CALL proc_add ( 'HelpDesk', 'Admins', 'AD' );
CALL proc_add ( 'Ali', 'Admins', 'AD' );
CALL proc_add (  'Ali', 'Users', 'AD');
CALL proc_add ( 'Burcu', 'Users', 'AD');
CALL proc_add ( 'Can', 'Users', 'AD' );
CALL proc_add (  'Managers', 'Users','AD' );
CALL proc_add ( 'Technicians', 'Users', 'AD');
CALL proc_add ('Demet', 'HelpDesk', 'AD' );
CALL proc_add ( 'Engin', 'HelpDesk', 'AD' );
CALL proc_add (  'Engin', 'Users', 'AD');
CALL proc_add ('Fuat', 'Managers', 'AD');
CALL proc_add ( 'G l', 'Managers', 'AD');
CALL proc_add ( 'Hakan', 'Technicians', 'AD');
CALL proc_add ('Irmak', 'Technicians', 'AD');
CALL proc_add ('ABCTechnicians', 'Technicians', 'AD');
CALL proc_add ('Jale', 'ABCTechnicians', 'AD');

```


### 删除（不可执行，参考逻辑）
```bigquery

DROP PROCEDURE proc_del;
#删除节点
delimiter $
CREATE PROCEDURE proc_del ( 
IN Id int)
BEGIN
	IF not EXISTS ( SELECT id FROM edge WHERE id=Id and hops = 0 ) THEN
     rollback;
	END IF;
	DROP TEMPORARY TABLE IF EXISTS purgeList;	
  DROP TEMPORARY TABLE IF EXISTS purgeList02;	
  DROP TEMPORARY TABLE IF EXISTS purgeList03;	
  DROP TEMPORARY TABLE IF EXISTS purgeList04;	
	create temporary table purgeList(
	  id int
	)Engine=InnoDB default charset utf8mb4;

insert into purgeList
  select id from edge where direct_edge_id = Id;

CREATE temporary TABLE purgeList02 like purgeList;
insert into purgeList02 select id from purgeList;
CREATE temporary TABLE purgeList03 like purgeList;
insert into purgeList03 select id from purgeList;
CREATE temporary TABLE purgeList04 like purgeList;
insert into purgeList04 select id from purgeList;

outer_label:BEGIN
while 1 = 1 do 
 insert into purgeList
	   select id from edge
		   where hops > 0
			   and (entry_edge_id in (select id from purgeList02) or exit_edge_id in (select id from purgeList03))
				 and id not in (select id from purgeList04);
	if ROW_COUNT() > 0 then
		leave outer_label;
	end if;
END WHILE;  
END outer_label;  
select * from purgeList;
-- 
-- delete from edge where id in(select id from purgeList);
END $
delimiter;
```