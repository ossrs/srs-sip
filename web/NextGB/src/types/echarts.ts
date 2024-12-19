import type { ComposeOption } from 'echarts/core'
import type { BarSeriesOption, LineSeriesOption } from 'echarts/charts'
import type { TitleComponentOption, TooltipComponentOption } from 'echarts/components'

export type ECOption = ComposeOption<
  BarSeriesOption | LineSeriesOption | TitleComponentOption | TooltipComponentOption
>
