# mongo-driver v2 변경점

- primitive.[타입] > bson.[타입] 형식으로 변경
- query 옵션 지정장식 변경

### prefix mgXXXX 으로 통일
- 중간에 대문자가 있으니까 개발할때 타이핑이 힘들다.
- 1글자 prefix 는 추후 혼동의 가능성이 높으므로 2글자로 통일