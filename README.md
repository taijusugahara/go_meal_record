# Meal Record
## 構成
<p><strong>React</strong> + <strong>Go</strong> + <strong>AWS</strong></p>
<p>(githubはGoの内容となります。Reactは<a href="https://github.com/taijusugahara/react_meal_record">こちら</a>になります。)</p>

## デプロイ
<li><strong>フロントエンド</strong>: AWS Amplifyにデプロイ</li>
<li><strong>バックエンド</strong>:  go,nginxをAWS ECSに。postgresをAWS RDSを使用してデプロイ</li>

## 使用技術
<p><strong>フロントエンド</strong></p>

<li>React(React Hooks, Redux Toolkit)</li>
<li>Typescript</li>
<li>MUI</li>
<li>CSS</li>
<li>HTML</li>
<br>
<p><strong>バックエンド</strong></p>

<li>Go(gin,gorm)</li>
<li>Postgres</li>
<li>Nginx</li>
<br>
<p><strong>インフラ?</strong></p>

<li>Docker(go,postgres,nginx)</li>
<li>AWS(ECR,ECS(Fargate),ALB,RDS,S3,Amplify etc)</li>
<li>CircleCI</li>

## 本アプリでやりたかったこと
<li>ReactとGoでSPAを作ること</li>
<li>コンテナ環境で開発して、デプロイもコンテナ環境で行うこと</li>

## 本アプリのコンセプト
<p>日々の食事内容を記録に残すこと。</p>
<p>私は「昨日、何食べたっけ？」と食事の内容を思い出せないことが多々あります。</p>
<p>そこで食事の度に、画像とテキストで内容を記録しておけば日々の食事を振り返ることができます。
また食事内容を振り返ることで新たな気づきがあります。例えば、「この料理は最近食べてないな。」や「この料理最近食べたのいつだっけ？」などです。</p>
<p>こうしたことにより本アプリは私たちの食事をより良くしてくれると考えます。</p>
