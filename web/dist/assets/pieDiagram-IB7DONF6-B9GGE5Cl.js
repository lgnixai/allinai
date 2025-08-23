import{p as V}from"./chunk-4BMEZGHF-ClXtSoD4.js";import{ag as S,a8 as z,aL as j,G as q,p as H,q as Z,c as J,g as K,e as Q,d as X,_ as p,m as G,y as Y,f as tt,H as et,M as at,U as rt,n as nt}from"./index-DLhw3Ufm.js";import{p as it}from"./radar-MK3ICKWK-D0mtwB3H.js";import{d as O}from"./arc-tXiW0gxU.js";import{o as st}from"./ordinal-Cboi1Yqb.js";import"./semi-ui-CsG5M4aO.js";import"./react-core-DtMAePju.js";import"./tools-2z7G0C69.js";import"./react-components-XvIrFW7f.js";import"./i18n-5-7yKGUP.js";import"./_baseUniq-Bs19QmAY.js";import"./_basePickBy-CLVtpQC8.js";import"./clone-BT4n6C1O.js";import"./init-Gi6I4Gst.js";function ot(t,a){return a<t?-1:a>t?1:a>=t?0:NaN}function lt(t){return t}function ct(){var t=lt,a=ot,m=null,o=S(0),u=S(z),x=S(0);function i(e){var r,l=(e=j(e)).length,g,A,h=0,c=new Array(l),n=new Array(l),v=+o.apply(this,arguments),w=Math.min(z,Math.max(-z,u.apply(this,arguments)-v)),f,T=Math.min(Math.abs(w)/l,x.apply(this,arguments)),$=T*(w<0?-1:1),d;for(r=0;r<l;++r)(d=n[c[r]=r]=+t(e[r],r,e))>0&&(h+=d);for(a!=null?c.sort(function(y,C){return a(n[y],n[C])}):m!=null&&c.sort(function(y,C){return m(e[y],e[C])}),r=0,A=h?(w-l*$)/h:0;r<l;++r,v=f)g=c[r],d=n[g],f=v+(d>0?d*A:0)+$,n[g]={data:e[g],index:r,value:d,startAngle:v,endAngle:f,padAngle:T};return n}return i.value=function(e){return arguments.length?(t=typeof e=="function"?e:S(+e),i):t},i.sortValues=function(e){return arguments.length?(a=e,m=null,i):a},i.sort=function(e){return arguments.length?(m=e,a=null,i):m},i.startAngle=function(e){return arguments.length?(o=typeof e=="function"?e:S(+e),i):o},i.endAngle=function(e){return arguments.length?(u=typeof e=="function"?e:S(+e),i):u},i.padAngle=function(e){return arguments.length?(x=typeof e=="function"?e:S(+e),i):x},i}var P=q.pie,F={sections:new Map,showData:!1,config:P},M=F.sections,W=F.showData,pt=structuredClone(P),ut=p(()=>structuredClone(pt),"getConfig"),gt=p(()=>{M=new Map,W=F.showData,Y()},"clear"),dt=p(({label:t,value:a})=>{M.has(t)||(M.set(t,a),G.debug(`added new section: ${t}, with value: ${a}`))},"addSection"),ft=p(()=>M,"getSections"),mt=p(t=>{W=t},"setShowData"),ht=p(()=>W,"getShowData"),R={getConfig:ut,clear:gt,setDiagramTitle:H,getDiagramTitle:Z,setAccTitle:J,getAccTitle:K,setAccDescription:Q,getAccDescription:X,addSection:dt,getSections:ft,setShowData:mt,getShowData:ht},vt=p((t,a)=>{V(t,a),a.setShowData(t.showData),t.sections.map(a.addSection)},"populateDb"),yt={parse:p(async t=>{const a=await it("pie",t);G.debug(a),vt(a,R)},"parse")},St=p(t=>`
  .pieCircle{
    stroke: ${t.pieStrokeColor};
    stroke-width : ${t.pieStrokeWidth};
    opacity : ${t.pieOpacity};
  }
  .pieOuterCircle{
    stroke: ${t.pieOuterStrokeColor};
    stroke-width: ${t.pieOuterStrokeWidth};
    fill: none;
  }
  .pieTitleText {
    text-anchor: middle;
    font-size: ${t.pieTitleTextSize};
    fill: ${t.pieTitleTextColor};
    font-family: ${t.fontFamily};
  }
  .slice {
    font-family: ${t.fontFamily};
    fill: ${t.pieSectionTextColor};
    font-size:${t.pieSectionTextSize};
    // fill: white;
  }
  .legend text {
    fill: ${t.pieLegendTextColor};
    font-family: ${t.fontFamily};
    font-size: ${t.pieLegendTextSize};
  }
`,"getStyles"),xt=St,At=p(t=>{const a=[...t.entries()].map(o=>({label:o[0],value:o[1]})).sort((o,u)=>u.value-o.value);return ct().value(o=>o.value)(a)},"createPieArcs"),wt=p((t,a,m,o)=>{G.debug(`rendering pie chart
`+t);const u=o.db,x=tt(),i=et(u.getConfig(),x.pie),e=40,r=18,l=4,g=450,A=g,h=at(a),c=h.append("g");c.attr("transform","translate("+A/2+","+g/2+")");const{themeVariables:n}=x;let[v]=rt(n.pieOuterStrokeWidth);v??(v=2);const w=i.textPosition,f=Math.min(A,g)/2-e,T=O().innerRadius(0).outerRadius(f),$=O().innerRadius(f*w).outerRadius(f*w);c.append("circle").attr("cx",0).attr("cy",0).attr("r",f+v/2).attr("class","pieOuterCircle");const d=u.getSections(),y=At(d),C=[n.pie1,n.pie2,n.pie3,n.pie4,n.pie5,n.pie6,n.pie7,n.pie8,n.pie9,n.pie10,n.pie11,n.pie12],D=st(C);c.selectAll("mySlices").data(y).enter().append("path").attr("d",T).attr("fill",s=>D(s.data.label)).attr("class","pieCircle");let L=0;d.forEach(s=>{L+=s}),c.selectAll("mySlices").data(y).enter().append("text").text(s=>(s.data.value/L*100).toFixed(0)+"%").attr("transform",s=>"translate("+$.centroid(s)+")").style("text-anchor","middle").attr("class","slice"),c.append("text").text(u.getDiagramTitle()).attr("x",0).attr("y",-400/2).attr("class","pieTitleText");const b=c.selectAll(".legend").data(D.domain()).enter().append("g").attr("class","legend").attr("transform",(s,k)=>{const E=r+l,_=E*D.domain().length/2,B=12*r,U=k*E-_;return"translate("+B+","+U+")"});b.append("rect").attr("width",r).attr("height",r).style("fill",D).style("stroke",D),b.data(y).append("text").attr("x",r+l).attr("y",r-l).text(s=>{const{label:k,value:E}=s.data;return u.getShowData()?`${k} [${E}]`:k});const I=Math.max(...b.selectAll("text").nodes().map(s=>(s==null?void 0:s.getBoundingClientRect().width)??0)),N=A+e+r+l+I;h.attr("viewBox",`0 0 ${N} ${g}`),nt(h,g,N,i.useMaxWidth)},"draw"),Ct={draw:wt},Pt={parser:yt,db:R,renderer:Ct,styles:xt};export{Pt as diagram};
